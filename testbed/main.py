#!/usr/bin/env python3

import subprocess
import time

def delete_nginx():
    out = None
    try:
        out = subprocess.check_output("kubectl delete -f ../deployments/pod.yml".split(" ")).decode()
    except:
        pass
    return out
    

def deploy_nginx():
    out = None
    try:
        out = subprocess.check_output("kubectl apply -f ../deployments/pod.yml".split(" ")).decode()
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

def test_no_interference(n, approach="tcc", subdir="no-interference"):
    for i in range(n):
        print('Testing no interference', i)
        print(delete_nginx())
        print(deploy_nginx())

        ip = wait_ip('web-deploy')

        save_csv(f'{subdir}/{approach}/iteration-{str(i).zfill(2)}.csv', start_hey(ip))

n = 30
test_no_interference(n, approach="k8s", subdir="1-thread")