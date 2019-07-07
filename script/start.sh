#!/bin/sh
#compile go build
go build .

#check if frolleague process is running
ps -ef | grep ./frolleague
if [ $? -eq 0 ]
then
  echo "frolleague Running..."
  echo "killing process..."
  ps -ef | grep ./frolleague | grep -v grep | awk '{print $2}' | xargs kill
  if [ $? -eq 0 ]
    then
    echo "frolleague process killed!"
  else
    echo "could not kill frolleague process" >&2
    exit 1
  fi
else
  echo "process not running"
fi

nohup ./frolleague &
if [ $? -eq 0 ]
then
  echo "frolleague process restarted"
  tail -f nohup.out
else
  echo "failed to start frolleague process" >&2
fi
