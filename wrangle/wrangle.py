import csv
import json 

file = open('dict/A.csv')

csvreader = csv.reader(file)

# Read all rows
rows = []
for row in csvreader:
        if row != []:
            rows.append(row)

# Convert rows to strings
strings = []
for row in rows:
    strings.append(row[0])

# Clean Data
for i in range(len(strings)):
    strings[i] = strings[i].replace(';', '')
    strings[i] = strings[i].replace(',', '')
    strings[i] = strings[i].replace('.', '')
    strings[i] = strings[i].replace('-', '')
    strings[i] = strings[i].replace('"', '')
    strings[i] = strings[i].replace('[', '')
    strings[i] = strings[i].replace(']', '')
    strings[i] = strings[i].replace(':', '')
    strings[i] = strings[i].replace('/', '')
    strings[i] = strings[i].replace('!', '')
    strings[i] = strings[i].replace('&', '')
    strings[i] = strings[i].replace('?', '')
    strings[i] = strings[i].replace('*', '')
    strings[i] = strings[i].replace('~', '')
    strings[i] = strings[i].replace('=', '')
    strings[i] = strings[i].replace('`', '')
    strings[i] = strings[i].replace('+', '')
    strings[i] = strings[i].replace('#', '')
    strings[i] = strings[i].replace('¡', '')
    strings[i] = strings[i].replace('–', '')
    strings[i] = strings[i].replace('^', '')
    strings[i] = strings[i].replace('$', '')
    strings[i] = strings[i].replace('{', '')
    strings[i] = strings[i].replace('}', '')
    strings[i] = strings[i].replace('\\', '')
    strings[i] = strings[i].replace('|', '')
    strings[i] = strings[i].replace('<', '')
    strings[i] = strings[i].replace('>', '')
    strings[i] = strings[i].replace('£', '')

# split word and defintion into dictionary
dict = {}
for i in range(len(strings)):
    for j in range(len(strings[i])):
        name, defn = strings[i].split(' ', 1)

        for k in range(len(defn)):
            if defn[k] == ')':
                defn = defn[k+1:]
                break
        
        #finish cleaning data
        defn = defn.replace('(', '')
        defn = defn.replace(')', '')

        dict[name] = defn.split()


# dump dictionary
with open("cleaned/A.json", "w") as outfile:
    json.dump(dict, outfile, indent=2)