#!/bin/bash

print_usage() {
  printf "Usage: RunTest.sh [OPTION] \n\
  
  OPTION:
    -h                      Show this help
    -p [Path to Profile]    Execute the test, running oscalkit code generate with the provided profile
    \n"
}

while getopts ':p:h' flag; do
  case "${flag}" in
    p) path=$OPTARG ;;
    h) print_usage ;
        exit ;;
    *) "Invalid option: -$OPTARG" ;
        print_usage ;
        exit ;;
  esac
done

if [ $OPTIND -eq 1 ]; then
    print_usage
    exit
fi

RED='\033[0;31m'
NC='\033[0m'
echo "Running test with profile path: $path"
go run cli/main.go generate code -p $path
if [ $? -eq 1 ]; then
    echo "${RED}Generate code failed. Exiting test.${NC}"
    rm output.go
    exit
fi
go run test_util/src/*.go -p $path
rm output.go