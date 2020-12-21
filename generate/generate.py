import json 
import random
import sys
import math

def print_usage():
	print("Usage: generate.py cities.json <maps.txt> <queries.txt> <number_queries> \n" +
	"\t cities.json = json file of city list\n" +
	"\t <maps.txt> = file name for empty txt file to put flights info\n" +
	"\t <queries.txt> = file name for empty txt file for customer requests\n" +
	"\t <number_queries> = total number of customer requests to put into queries.txt file\n" +
	"Sample Runs:\n" +
	"\t./generate cities.json maps250.txt queries250.txt 250 -- Generates 250 queries \n")


if len(sys.argv) == 5:
	city_map = sys.argv[1]
	map_filepath = sys.argv[2]
	query_filepath = sys.argv[3]
	num_queries = int(sys.argv[4])
else:
	print_usage()

# Opening JSON file 
with open(city_map) as json_file: 
    data = json.load(json_file) 

# Get comprehensive list of cities 
cities = []   

for entry in data:
	name = entry["name"] + ", " + entry["country"]
	cities.append(name)

# Get subset of cities needed to get number of unique queries 
cities_needed = math.ceil(math.sqrt(num_queries)) + 10
print(cities_needed)

city_sublist = random.sample(cities, cities_needed)

# Create maps.txt file 
maps = []

j = 1
for source in city_sublist:
	for destination in city_sublist:
		if source != destination:
			price = random.randint(100, 2000)
			flight = '{"origin": "' + source + '", "destination": "' + destination + '", "price": ' + str(price) + '}'
			maps.append(flight)
			print(flight)
			print(j)
			j += 1

with open(map_filepath, 'w') as map_file:
    map_file.write('\n'.join(maps))
map_file.close()

# Create queries.txt file 
queries = []

i = 1
for source in city_sublist:
	for destination in city_sublist:
		if i <= num_queries: 
			query = '{"id": ' + str(i) + ', "origin": "' + source + '", "destination": "' + destination + '"}'
			print(query)
			print(i)
			queries.append(query)
			i += 1
		else:
			break

with open(query_filepath, 'w') as query_file:
    query_file.write('\n'.join(queries))
query_file.close()

