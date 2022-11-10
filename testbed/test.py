#!/usr/bin/env python3 

from main import get_pod_metrics
from pprint import pprint
from time import time

x = get_pod_metrics('web-deploy', 'cpu')
print(x)

