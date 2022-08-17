import csv
import json
from nltk.stem import WordNetLemmatizer

lemma = WordNetLemmatizer()

lem = False

for i in range(ord('A'), ord('Z')+1):

    rows = []
    strings = []
    
    char = chr(i)

    print(char)

    file = open('dict/'+char+'.csv')

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

        strings[i] = strings[i].lower()

    print("data cleaned")

    # split word and defintion into dictionary
    dict = {}
    for i in range(len(strings)):
        # split name and definiton
        name, defn = strings[i].split(' ', 1)

        # remove word type from definiton
        for k in range(len(defn)):
            if defn[k] == ')':
                defn = defn[k+1:]
                break
        
        # finish cleaning data
        defn = defn.replace('(', '')
        defn = defn.replace(')', '')

        # split definiton into list
        defn = defn.split()

        # do stemming on words
        if (lem):
            #name = lemma.lemmatize(name)
            for i in range(len(defn)):
                defn[i] = lemma.lemmatize(defn[i])

        # save word/def to dictionary 
        dict[name] = defn


    # dump dictionary
    with open("cleaned/"+char+".json", "w") as outfile:
        json.dump(dict, outfile, indent=2)