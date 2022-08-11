import csv
import json 
from nltk.stem import PorterStemmer

porter = PorterStemmer()

stem = True

rows = []
strings = []

for i in range(ord('A'), ord('Z')+1):
    file = open('dict/'+chr(i)+'.csv')

    csvreader = csv.reader(file)

    # Read all rows
    for row in csvreader:
            if row != []:
                rows.append(row)

# Convert rows to strings
for row in rows:
    strings.append(row[0])

print("data read")

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

print("data cleaned")

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

        # split definiton into list
        defn = defn.split()

        # do stemming on words
        if (stem):
            name = porter.stem(name)
            for i in range(len(defn)):
                defn[i] = porter.stem(defn[i])

        dict[name] = defn


# dump dictionary
fn = "test.json"
with open("cleaned/"+fn, "w") as outfile:
    json.dump(dict, outfile, indent=2)