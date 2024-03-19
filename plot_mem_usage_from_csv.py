#!/usr/bin/env python3

import matplotlib.pyplot as plt
import pandas as pd
import numpy as np
import os
import sys

mem_logs_suffix = "_memory_usage.csv"


def plot_memory_usage_from_csv(folder_path, test_name):
    if not os.path.isdir(folder_path):
        print(f"Path: '{folder_path}' doesn't exist or is not a directory")
        exit(2)
    files = os.listdir(folder_path)
    memory_logs = [f for f in files if f.endswith(mem_logs_suffix)]
    if len(memory_logs) == 0:
        print("Didn't find any memory logs for the given test name. exiting..")
        exit(0)

    average_mem_use = []

    alg_names = [file_name[:-len(mem_logs_suffix)] for file_name in memory_logs]
    longest_alg_name_size = max(len(alg_name) for alg_name in alg_names)
    # print()
    # print(f"\tFound {len(memory_logs)} memory log files:   {', '.join(memory_logs)}")

    plt.rcParams['axes.xmargin'] = 0
    plt.rcParams['axes.ymargin'] = 0

    output_file_paths = []
    font_size = {'size': 24}
    # Create the plot
    scale = 2
    f, axes = plt.subplots(2, len(memory_logs), figsize=(16 * scale, 9 * scale))
    f.suptitle(f"{test_name}, Results:", fontsize=28, fontweight='bold')

    for idx, memlog_file_name in enumerate(memory_logs):

        alg_name = memlog_file_name[:-len(mem_logs_suffix)]
        df = pd.read_csv(folder_path + "/" + memlog_file_name)

        # Convert the time and memory columns to NumPy arrays
        time_nanosecs = np.array(df['Time'])
        memory = np.array(df['Memory (MB)'])


        time_milisecs = [t/ (10**6) for t in time_nanosecs]

        avg_memory_size = 0 if len(memory) == 0 else sum(memory)/len(memory)


        # print(f"printing time and mem for file: {memlog_file_name}")
        # print()
        # for time_mili, mem in zip(time_milisecs, memory):
        #     print(f"Time: {time_mili}, Memory: {mem}")

        # first plot
        top_plot = axes[0, idx]
        top_plot.plot(time_milisecs, memory, label='Memory Usage vs Time')
        top_plot.set_xlabel('Time(Miliseconds)', fontdict=font_size)
        top_plot.set_ylabel('Memory (MB)', fontdict=font_size)
        top_plot.set_title(f"Memory Usage vs Time - {alg_name}",fontdict=font_size)
        top_plot.tick_params(axis='both', which='major', labelsize=20)
        top_plot.legend(fontsize=20)

        # second plot
        bottom_plot = axes[1, idx]
        data_frequencies, x_bins, patches = bottom_plot.hist(x=memory, bins='auto', color='#0504aa',alpha=0.7, rwidth=0.7)
        bottom_plot.grid(axis='y', alpha=0.75)
        bottom_plot.tick_params(axis='both', which='major', labelsize=20)
        bottom_plot.set_xlabel("Memory (MB)", fontdict=font_size)
        bottom_plot.set_ylabel('Frequency', fontdict=font_size)
        bottom_plot.set_title(f"Memory Usage Distribution - {alg_name}", fontdict=font_size)
        bottom_plot.text(max(x_bins) * 0.8, max(data_frequencies) * 0.6, f"median={np.median(memory):.3f}\nmean={np.mean(memory):.3f}\nstd={np.std(memory):.3f}", fontdict=font_size)
        maxfreq = data_frequencies.max()
        bottom_plot.set_ylim(ymax=np.ceil(maxfreq / 10) * 10 if maxfreq % 10 else maxfreq + 10)



        plt.subplots_adjust(hspace=0.4, bottom=0.1, top=0.9, left=0.1, right=0.97)
        # output_file_path = folder_path + "/" + test_name + "_" + alg_name + "_memory_usage_plot.png"
        # plt.savefig(output_file_path)
        # plt.close(f)

        # padding = " " * (longest_alg_name_size - len(alg_name))
        # print(f"{idx}: Plot for {alg_name}{padding} ---------------> {output_file_path}")
        # output_file_paths.append(output_file_path)
        # print(f"\t\tAverage memory consumption during test for {alg_name}: {avg_memory_size} MB")
        # print()

        average_mem_use.append(avg_memory_size)

    print()
    print(f"The difference in average memory consumption was {max(average_mem_use) - min(average_mem_use):.3f} MB")
    print()

    output_file_path = folder_path + "/" + test_name + "_memory_usage_plot.png"
    plt.savefig(output_file_path)
    plt.close(f)

    print(f"plot saved in ---------------> {output_file_path}")


if __name__ == '__main__':

    if len(sys.argv) != 3:
        print(f"{sys.argv[0]}: Wrong number of arguments. expected 2, got {len(sys.argv) - 1}.")
        print("\tparam #1: test name - the name of the that we are plotting")
        print("\tparam $2: target folder - path to the folder in which the memory logs reside. the plot will be saved in the same dir.")
        exit(1)

    # target_dir = os.path.dirname(__file__) + "/server_test_results"
    test_name = sys.argv[1]
    target_dir = sys.argv[2]


    plot_memory_usage_from_csv(target_dir, test_name)
    

