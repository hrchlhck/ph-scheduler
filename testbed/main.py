#!/usr/bin/env python3

import subprocess
import time
import traceback

from operator import add
from functools import reduce
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

def get_pod_ip(pod_name: str):
    out = subprocess.check_output("kubectl get pods -o wide --no-headers".split(" ")).decode()

    def apply(x):
        x = x.split('\n')
        x = [o.split() for o in x if o]
        x = [p[5] for p in x if p[0] == pod_name]

        return x[0] != '<none>', x[0]

    return apply(out)

def save_csv(name: str, data: str) -> None:
    print('Saving data at', name)
    with open(name ,'w') as fp:
        fp.write(data)

def wait_ip(pod_name: str):
    exists, ip = get_pod_ip(pod_name)

    while not exists:
        exists, ip = get_pod_ip(pod_name)
        time.sleep(0.5)
        print('Waiting ip')
    
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

def test_interference(n, threads: int, approach="tcc", subdir="4-thread", policy="bestfit"):
    s = start_scheduler(approach, policy)

    for i in range(n):
        print('Testing interference', i, subdir)
        print(delete_interference(threads))
        print(create_interference(threads))
        print(delete_nginx(approach))
        print(deploy_nginx(approach))

        ip = wait_ip('web-deploy')

        save_csv(f'{subdir}/{policy}/{approach}/iteration-{str(i).zfill(2)}.csv', start_hey(ip))

    s.terminate()

n = 30
subdir = "4-thread"
approaches = ["tcc", "k8s"]
policies = ["bestfit", "worstfit", "firstfit"]

for a in approaches:
    for p in policies:
        test_interference(n, 4, approach=a, subdir=subdir, policy=p)
