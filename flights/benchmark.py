import math
import pandas as pd
import seaborn as sns
import subprocess

Benchmarks = {
		"title": "Flights Search Engine Speedup Graph",
		"x_label": "Number of Threads (N)",
		"y_label": "Speedup",
		"threads": [
			1,
			2,
			4,
			6,
			8
		],
		"filename": "flights.go",
		"timing_output": "timing",
		"lines": [
			{
				"Data Size": "25000",
				"map_file": "maps25000.txt",
				"query_file": "queries25000.txt",
				"seq": True,
				"par": True
			},
			{
				"Data Size": "50000",
				"map_file": "maps50000.txt",
				"query_file": "queries50000.txt",
				"seq": True,
				"par": True
			},
			{
				"Data Size": "75000",
				"map_file": "maps75000.txt",
				"query_file": "queries75000.txt",
				"seq": True,
				"par": True
			},
			{
				"Data Size": "100000",
				"map_file": "maps100000.txt",
				"query_file": "queries100000.txt",
				"seq": True,
				"par": True
			},
		],
		"speedup": False,
		"repeat": 5,
		"output_file": "speedup.png"
}

title = Benchmarks["title"]
x_label = Benchmarks["x_label"]
y_label = Benchmarks["y_label"]
filename = Benchmarks["filename"]
threads = Benchmarks["threads"]
lines = Benchmarks["lines"]
speedup = Benchmarks["speedup"]
repeat = int(Benchmarks["repeat"])
output_file = Benchmarks["output_file"]
speedup = Benchmarks["speedup"]
timing_output = Benchmarks["timing_output"]

#Given a time file, retrieve the real time
#Format: '0m2.234s'
def getTime(file):
	f = open(file, "r")
	lines = []
	for line in f:
	  lines.append(line.split())

	return lines[1][1]

#Given time (string), get seconds (float)
#Arg: '0m2.234s' 
def timeToSecs(time):
	split = time.split("m")
	minute = split[0]
	minToSecs = float(minute) * 60
	secString = float(split[1][:-1])
	return minToSecs + secString

#Get average secs (float) given list of times (string) 
def getAverage(timeArray):
	sums = 0
	for time in timeArray:
		secs = timeToSecs(time)
		sums += secs
	return sums/ len(timeArray)

#Calculate speedup
def Speedup(serialSecs, parallelSecs):
	return serialSecs / parallelSecs

#Returns an output graph given data of format
#data = [['group', threads (int), secs(int)], [], []]
def speedupGraph(data):
    df = pd.DataFrame(data, columns=['Data Size', x_label, y_label])
    print(df)
    fig = sns.lineplot(data=df, x=x_label, y=y_label, hue='Data Size').set_title(title)
    fig.figure.savefig(output_file)


data = []
for line in lines: 
	base = 0
	#Sequential
	if line["seq"] == True: 
		#seq_cmd = "time(go run " + filename + ") < " + line["input_file"] + " > out.txt 2> " + timing_output
		times = []
		for i in range(1, repeat+1):
			timing_output_file = timing_output + "_" + str(line["Data Size"]) + "_seq_" + str(i)
			seq_cmd = "time(go run " + filename + " " + line["map_file"] + " " + line["query_file"] + ") > out.txt 2> " + timing_output_file
			print(seq_cmd)
			subprocess.call(["bash", "-c", seq_cmd])
			time = getTime(timing_output_file)
			#time = getTime(timing_output)
			times.append(time)
		print("Sequential Times: ") 
		print(times)
		base = getAverage(times)
		print("Average Sequential Time " + str(base))

	#Parallel
	if line["par"] == True:
		for thread in threads:
			P = line["Data Size"]
			N = thread
			#par_cmd = "time(go run " + filename + " " + str(N) + " " + str(B) + ") < " + line["input_file"] + " > out.txt 2> " + timing_output
			times = []
			for i in range(1, repeat+1):
				timing_output_file = timing_output + "_" + str(line["Data Size"]) + "_par_" + str(thread) + "_" + str(i)
				par_cmd = "time(go run " + filename + " " + str(N) + " " + line["map_file"] + " " + line["query_file"] + ") > out.txt 2> " + timing_output_file
				print(par_cmd)
				subprocess.call(["bash", "-c", par_cmd])
				time = getTime(timing_output_file)
				#time = getTime(timing_output)
				times.append(time)
			print("Running lines for thread: " + str(thread))
			print("Parallel Times: ")
			print(times)
			avgTime = getAverage(times)
			print("Average Parallel Time " + str(avgTime))
			speedup = Speedup(base, avgTime)
			summary = [P, thread, speedup]
			print("Lines, Thread, Speedup: ")
			print(summary)
			data.append(summary)
print(data)	
speedupGraph(data)









