#!/bin/bash

BASE_DIR=memory

for interference in $(ls $BASE_DIR); do
		_dd=$BASE_DIR/$interference
		for policy in $(ls $_dd); do 
				_ddd=$BASE_DIR/$interference/$policy
				for ap in $(ls $_ddd); do
						d=$BASE_DIR/$interference/$policy/$ap
						n=$(ls $d | wc -l)
						if [[ $(ls $d | wc -l) -eq "30" ]]; then
								echo $d ok 
						else
								echo $d fail
						fi
				done
		done
done

