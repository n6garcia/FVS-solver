from nltk.wsd import lesk
from nltk.tokenize import word_tokenize
from nltk.corpus import wordnet


c = lesk(word_tokenize('face of awf earth'), "awf")
if c != None:
    print(c, c.definition())
    #print(c.name())
    #syns = wordnet.synset(c.name())
    #print(syns.definition())
    #print(word_tokenize('face of. awf-wf earth'))


allWords = [n for n in wordnet.all_synsets()]
print(len(allWords))
print(allWords[0].name())