#!/bin/bash
# 运行 20 次循环
for i in {1..20000}
do
   curl http://localhost/ping
   echo " - Request $i completed"
   sleep 0.5
done