#!/usr/bin/env python3

import requests
import subprocess
import time
import traceback
import re

from datetime import datetime as dt
from pandas import DataFrame
from collect import Collector
from pathlib import Path

def delete_nginx(approach: str):
    out = None
    try:
        out = subprocess.check_output(f"kubectl delete -f ../deployments/pod_{approach}.yml".split(" ")).decode()
    except:
        pass
    return out


def deploy_nginx(approach: str):
    out = None
    try:
        out = subprocess.check_output(f"kubectl apply -f ../deployments/pod_{approach}.yml".split(" ")).decode()
    except:
        pass
    return out

def start_hey(pod_ip: str):
    out = None
    try:
        print('Started hey')
        out = subprocess.check_output(f"docker run -it --rm vpemfh7/hey -c 50 -n 30000 -o csv -z 15s http://{pod_ip}".split(" ")).decode()
    except:
        pass
    return out

def start_hey_detach(pod_ip: str):
    out = None
    try:
        print('Started hey')
        out = subprocess.check_output(f"docker run -it --rm -d vpemfh7/hey -c 50 -n 30000 -o csv -z 15s http://{pod_ip}".split(" ")).decode()
    except:
        pass
    return out

def get_pod_ip(pod_name: str):
    out = subprocess.check_output("kubectl get pods -o wide --no-headers".split(" ")).decode()

    def apply(x):
        x = x.split('\n')
        x = [o.split() for o in x if o]
        x = [p[5] for p in x if p[0] == pod_name]

        return x[0] != '<none>', x[0]

    return apply(out)

def get_pod_node(pod_name: str):
    out = subprocess.check_output("kubectl get pods -o wide --no-headers".split(" ")).decode()

    def apply(x):
        x = x.split('\n')
        x = [o.split() for o in x if o]
        x = [p[6] for p in x if p[0] == pod_name]

        return x[0] != '<none>', x[0]

    return apply(out)

def get_metrics_server_ip(pod_name: str, pod_node: str) -> tuple:
    out = subprocess.check_output("kubectl get pods -o wide --no-headers".split(" ")).decode()
    
    out = out.split('\n')
    out = [o for o in out if 'metrics-server' in o and pod_node in o]
    out = re.findall(r'(([0-9]{1,4})\.([0-9]{1,4})\.([0-9]{1,4})\.([0-9]{1,4}))', '\n'.join(out))
    return out[0][0]

def filter_containers(data: dict) -> dict:
    return {k:v for k, v in data.items() if 'POD' not in k}

def get_node_containers(ms_ip: str) -> dict:
    return requests.get("http://" + ms_ip + "/docker").json()['containers']

def wait_container_metrics_server(pod_name: str) -> None:

    _, pod_node = get_pod_node(pod_name)
    ms_ip = get_metrics_server_ip(pod_name, pod_node)

    req = get_node_containers(ms_ip)
    req = filter_containers(req)
    i = 0    
    cond = any(i for i in req if pod_name in i)
    while not cond:
        _, pod_node = get_pod_node(pod_name)
        ms_ip = get_metrics_server_ip(pod_name, pod_node)

        req = get_node_containers(ms_ip)
        req = filter_containers(req)
        cond = any(i for i in req if pod_name in i)
        print(f'[ {i}s ] Waiting MS ip')
        time.sleep(1)
        i += 1

    return ms_ip    

def get_pod_id(pod_name: str):
    wait_ip(pod_name)

    ms_ip = wait_container_metrics_server(pod_name)

    req = requests.get("http://" + ms_ip + "/docker").json()['containers']
    
    key = None

    for k in req:
        if pod_name in k and 'POD' not in k:
            key = k

    return req[key]

def get_pod_metrics(pod_name: str, device: str) -> DataFrame:
    pod_id = get_pod_id(pod_name)
    _, pod_node = get_pod_node(pod_name)
    ms_ip = get_metrics_server_ip(pod_name, pod_node)

    c = Collector('http://' + ms_ip + f"/docker/{pod_id}", device)

    data = c.collect()

    return c.stats(data)

def save_csv(name: str, data: str) -> None:
    print('Saving data at', name)
    with open(name ,'w') as fp:
        fp.write(data)

def wait_ip(pod_name: str):
    exists, ip = get_pod_ip(pod_name)

    i = 0
    while not exists:
        exists, ip = get_pod_ip(pod_name)
        time.sleep(1)
        print(f'[ {i}s ] Waiting ip')
        i+=1
    
    return ip

def start_scheduler(scheduler_name: str, policy: str) -> subprocess.Popen:
    proc = None
    try:
        print(f'Initializing scheduler "{scheduler_name}" with policy "{policy}"')
        proc = subprocess.Popen(f"{Path('.').cwd().parent / 'scheduler'} {scheduler_name} {policy}".split(" "), stdout=subprocess.DEVNULL, stderr=subprocess.STDOUT)
    except:
        traceback.print_exc()
    return proc

def test_no_interference(n, approach="tcc", subdir="no-interference", policy="bestfit"):
    s = start_scheduler(approach, policy)

    for i in range(n):
        print('Testing no interference', i)
        print(delete_nginx(approach))
        print(deploy_nginx(approach))

        ip = wait_ip('web-deploy')

        save_csv(f'{subdir}/{policy}/{approach}/iteration-{str(i).zfill(2)}.csv', start_hey(ip))

    s.terminate()

def wait_deletion():
    done = False
    i = 0
    while not done:
        out = subprocess.check_output("kubectl get pods -o wide --no-headers".split(" ")).decode()
        out = out.split('\n')
        out = [i.split()[0] for i in out if len(i) > 0]
        out = [i for i in out if 'sysbench' in i]
        out = ', '.join(out)
    
        if len(out) == 0:
            done = True

        time.sleep(1)
        print(f'[ {i}s ]Waiting for pods to be deleted')
        i += 1

def create_interference(n_threads=4) -> None:
    out = None
    try:
        out = subprocess.check_output(f"kubectl apply -f ../deployments/sysbench_{n_threads}threads.yml".split(" ")).decode()
    except:
        pass
    return out

def delete_interference(n_threads=4) -> None:
    out = None
    try:
        out = subprocess.check_output(f"kubectl delete -f ../deployments/sysbench_{n_threads}threads.yml".split(" ")).decode()
        wait_deletion()
    except:
        pass
    return out


def create_interference_memory(mem_block_size=4) -> None:
    out = None
    try:
        out = subprocess.check_output(f"kubectl apply -f ../deployments/sysbench_{mem_block_size}.yml".split(" ")).decode()
    except:
        pass
    return out

def delete_interference_memory(mem_block_size=4) -> None:
    out = None
    try:
        out = subprocess.check_output(f"kubectl delete -f ../deployments/sysbench_{mem_block_size}.yml".split(" ")).decode()
        wait_deletion()
    except:
        pass
    return out

def test_interference(n, threads: int, device='memory', approach="tcc", subdir="4-thread", policy="bestfit"):
    for i in range(n):
        print(f'[ {i} ] Testing interference', subdir)
        
        if subdir != 'no-interference':
            print(f'[ {i} ] No interference')
            print(delete_interference_memory(threads))

        s = start_scheduler(approach, policy)
        
        if subdir != 'no-interference':
            print(f'[ {i} ] No interference')
            print(create_interference_memory(threads))
        
        print(delete_nginx(approach))
        print(deploy_nginx(approach))

        ip = wait_ip('web-deploy')

        start_hey_detach(ip)

        # save_csv(f'{subdir}/{policy}/{approach}/iteration-{str(i).zfill(2)}.csv', start_hey(ip))
        data = get_pod_metrics('web-deploy', device)
        data.to_csv(f'{device}/{subdir}/{policy}/{approach}/iteration-{str(i).zfill(2)}.csv', index=True)

        s.terminate()

if __name__ == '__main__':
    n = 30
    device = 'memory'
    subdirs = ['no-interference', '4gb', '8gb']
    approaches = ["tcc", "k8s"]
    policies = ["bestfit", "worstfit", "firstfit"]
    
    st = dt.now()
    for s in subdirs[1:]:
        for a in approaches:
            for p in policies:
                
                if s == subdirs[0]:
                    memory_block_size = 0
                else:
                    memory_block_size = int(s[0])
                test_interference(n, s, device, approach=a, subdir=s, policy=p)
    et  = dt.now()
    print(f'Testbed took {et-st}')