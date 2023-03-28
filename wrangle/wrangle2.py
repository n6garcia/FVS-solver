from nltk.wsd import lesk
from nltk.tokenize import word_tokenize
from nltk.corpus import wordnet
import json

dict = {}

allWords = [n for n in wordnet.all_synsets()]

for word in allWords:
    ID = word.name()
    name = ID.split('.')[0]
    if '_' in name:
        continue
    origDef = word.definition()
    regexDef = word.definition()
    regexWords = []
    Mappings = []

    tkns = word_tokenize(origDef)
    regTkns = word_tokenize(origDef)

    for i in range(len(tkns)):
        tkn = tkns[i]
        c = lesk(origDef, tkn)
        if c is None:
            continue
        elif '_' in c.name().split('.')[0]:
            continue
        elif tkn != c.name().split('.')[0]:
            continue
        else:
            regTkns[i] = "%s"
            regexWords.append(tkn)
            Mappings.append(c.name())

    regexDef = ' '.join(regTkns)

    tpl = (name,origDef,regexDef,regexWords,Mappings)

    dict[ID] = tpl


with open("wordnet/wn.json", "w") as outfile:
    json.dump(dict, outfile, indent=2)