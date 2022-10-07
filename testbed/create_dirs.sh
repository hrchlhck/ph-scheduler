#!/bin/bash

TESTS=(no-interference 1-thread 4-thread 8-thread)
APPROACHES=(tcc k8s)
POLICIES=(bestfit worstfit firstfit)

for t in ${TESTS[@]}; do
		for p in ${POLICIES[@]}; do
			for a in ${APPROACHES[@]}; do
					mkdir -p $t/$p/$a;
			done
		done
done
