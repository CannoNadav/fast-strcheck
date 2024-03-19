#!/bin/bash

client_test_name="$1"

full_client_test_name="Test${client_test_name}Client"
valid_client_test_names=(HighLoadTest ResponseValidatingTest)
script_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
test_results_dir="${script_dir}/server_test_results"

if [ -z "$client_test_name" ]; then
	echo "positional argument #1 missing: TEST_NAME"
	echo did not select a test to run. options are: "${valid_client_test_names[@]}"
	echo exiting
	exit 1
elif [ "$client_test_name" != "${valid_client_test_names[0]}" ] && [ "$client_test_name" != "${valid_client_test_names[1]}" ]; then
	echo ilegal value for client test name: "$client_test_name". valid options are: "${valid_client_test_names[@]}"
	exit 1
fi


ds_algorithms=(MapWordsCatalog MemBlockWordsCatalog)

echo "cleaning results dir(if exists): ${script_dir}/server_test_results"
if [ -d "${test_results_dir}" ]; then
	rm -f "${test_results_dir}"/*.csv
	rm -f "${test_results_dir}"/*.pprof
	rm -f "${test_results_dir}"/*.log
	rm -f "${test_results_dir}"/*.png
fi

mkdir -p "${test_results_dir}"
echo
echo
echo "running test: ${client_test_name} with both data structure algorithms: ${ds_algorithms[*]}"
echo

go clean

for algorithm in "${ds_algorithms[@]}"; do

	echo "    ---------->  starting ${client_test_name} with algorithm '${algorithm}'"

	go test ./server/ -count=1 -v -run TestServerWithProfiling -WordsCatalogAlg="${algorithm}" &> "${test_results_dir}/ServerProfilingLogs_${algorithm}.log" &
	server_pid=$!

	go test ./test_client -count=1 -v -run "${full_client_test_name}" &> "./server_test_results/${client_test_name}_${algorithm}.log"  &
	client_pid=$!

	# echo server pid: $server_pid client pid: $client_pid

	server_finished=false
	client_finished=false

	while [[ "${server_finished}" != true ]] && [[ "${client_finished}" != true ]]; do
	    # Check if server process is still running
	    # echo checking server
	    if ! ps -p $server_pid &> /dev/null; then
	    	# echo server is down
	        wait $server_pid # get exit status of server
	        server_exit_code=$?
	        # echo server exit code: $server_exit_code
	        if [ $server_exit_code -eq 0 ]; then
	        	server_finished=true
	        else
	        	echo "Server terminated with error, see logs in ${test_results_dir} for details"
	        	echo "killing client with pid=${client_pid} to terminate"
	        	kill -9 "$client_pid"
	        	exit 2
	        fi
	    fi

	    # echo checking client
	    # Check if client process is still running
	    if ! ps -p $client_pid &> /dev/null; then
	    	# echo client is down
	        wait $client_pid # get exit status of client
	        client_exit_code=$?
	        if [ $client_exit_code -eq 0 ]; then
	        	client_finished=true
	        else
	        	echo "Client terminated with error, see logs in ${test_results_dir} for details"
	        	echo "killing server with pid=${server_pid} to terminate"
	        	kill -9 "$server_pid"
	        	exit 2
	        fi
	    fi

	    sleep 0.5
	done


done

echo
echo  Tests finished, printing server side summary of tests:
echo

did_fail=false
server_mapcatalog_log="${test_results_dir}/ServerProfilingLogs_MapWordsCatalog.log"

server_blockcatalog_log="${test_results_dir}/ServerProfilingLogs_MemBlockWordsCatalog.log"

if grep -q "FAIL" "${server_mapcatalog_log}"; then
  cat "${server_mapcatalog_log}"
  did_fail=true
else
  grep -A 5 "Results(Summary):" "${server_mapcatalog_log}"
fi

echo 

if grep -q "FAIL" "${server_blockcatalog_log}"; then
  cat "${server_blockcatalog_log}"
  did_fail=true
else
	grep -A 5 "Results(Summary):" "${server_blockcatalog_log}"
fi

echo
echo  Full logs from both client and server can be found in the '"server_test_results"' directory

if [ "${did_fail}" = false ]; then
	echo  plotting memory consumption for both algorithms:
	echo

	chmod u+x ./plot_mem_usage_from_csv.py
	if ./plot_mem_usage_from_csv.py  "${client_test_name}" "${test_results_dir}"; then 
		# display the plots
		xdg-open "${test_results_dir}/${client_test_name}_memory_usage_plot.png"
	fi
fi