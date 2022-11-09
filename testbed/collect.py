import requests
import pandas as pd

from time import sleep as t_sleep
from time import time

def subtract_dicts(d1: dict, d2: dict) -> dict:
    if len(d1) != len(d2):
        raise ValueError("Dictionaries has different lenghts")
    
    if d1.keys() != d2.keys():
        raise ValueError("Dictionaries has different keys")
    
    ret = dict()

    for k in d1:
        ret[k] = d2[k] - d1[k]
    
    return ret

class Collector:
    def __init__(self, mc_ip: str, device: str, interval=0.01, total_time=15):
        self._ip = mc_ip
        self._interval = interval
        self._total_time = total_time
        self._device = device
    
    def stats(self, data: list) -> dict:
        return pd.DataFrame(data=data).describe().T

    def collect(self) -> list:
        ret = list()
        st = time()
        while time() - st < self._total_time:
            old = requests.get(self._ip + "/" + self._device).json()
            t_sleep(self._interval)
            new = requests.get(self._ip + "/" + self._device).json()
            new = subtract_dicts(old, new)

            ret.append(new)
        return ret

    def _save(self, output_name: str) -> None:
        pass