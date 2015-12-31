import os

try:
    os.remove('pairing.csv')
    os.remove('pairing_foods.txt')
except:
    pass

pairs = {}
with open('srep00196-s2.csv') as f:
    for line in f:
        line = line.replace('olive_oil','olivoil')
        line = line.replace('uncured_','')
        line = line.replace('cured_','')
        line = line.replace('roasted_','')
        line = line.replace('_oil','')
        line = line.replace('_juice','')
        line = line.replace('baked_','')
        line = line.replace('grilled_','')
        line = line.replace('fried_','')
        line = line.replace('raw_','')
        line = line.replace('smoked_','')
        line = line.replace('mashed_','')
        line = line.replace('cooked_','')
        line = line.replace('unprocessed_','')
        line = line.replace('bitter_','')
        line = line.replace('water_','')
        line = line.replace('black_','')
        line = line.replace('boiled_','')
        line = line.replace('dried_','')
        line = line.replace('crips_','')
        line = line.replace('summer_','')
        line = line.replace('lean_','')
        line = line.replace('red_','')
        if "vegetable" not in line and "leaf" not in line and "#" not in line and "meat" not in line and "tea," not in line and "clove" not in line and "pork" not in line and "beef" not in line and "chicken" not in line and "fish" not in line and line[0:5]!="bean,":
            line = line.replace("_"," ")
            foods = line.split(',')
            food1 = foods[0]
            food2 = foods[1]
            score = int(foods[2])
            if score > 1:
                if foods[1] < foods[0]:
                    food1 = foods[1]
                    food2 = foods[0]
                if food1 in pairs:
                    if food2 in pairs[food1]:
                        pairs[food1][food2] += score
                    else:
                        pairs[food1][food2] = score
                else:
                    pairs[food1] = {}
                    pairs[food1][food2] = score



import numpy
nums = []
for food1 in pairs:
    for food2 in pairs[food1]:
        nums.append(pairs[food1][food2])

nums = sorted(nums)
cdf = numpy.cumsum(nums,dtype=float)/numpy.sum(nums,dtype=float)
lastNum = nums[0]
cdfTable = {}
for i in range(len(nums)):
    if nums[i] != lastNum:
        cdfTable[int(nums[i-1])] = int(100*cdf[i-1])
cdfTable[int(nums[i])] = int(100*cdf[i])

uniqueFoods = []
with open("pairing.csv","w") as f:
    for food1 in pairs:
        for food2 in pairs[food1]:
            f.write(food1 + "," + food2 + "," + str(cdfTable[pairs[food1][food2]]) + "\n")
            if food1 not in uniqueFoods:
                uniqueFoods.append(food1)
            if food2 not in uniqueFoods:
                uniqueFoods.append(food2)

with open('pairing_foods.txt','w') as f:
    for food in uniqueFoods:
        f.write(food + "\n")



import random

uniqueFoods = []
with open('pairing_foods.txt','r') as f:
    for line in f:
        uniqueFoods.append(line.strip())

pairs = {}
with open("pairing.csv","r") as f:
    for line in f:
        foods = line.split(',')
        food1 = foods[0]
        food2 = foods[1]
        score = int(foods[2])
        if food1 not in pairs:
            pairs[food1] = {}  
        if food2 not in pairs:
            pairs[food2] = {}
        pairs[food1][food2] = score
        pairs[food2][food1] = score


bootstraps = []

for ijkl in range(2000):
    ingredients = random.sample(uniqueFoods,random.randint(6,15))
    score = 0
    num = 0
    for i in range(len(ingredients)):
        for j in range(1,len(ingredients),1):
            if j > i and ingredients[i] in pairs:
                if ingredients[j] in pairs[ingredients[i]]:
                    print(ingredients[i],ingredients[j],pairs[ingredients[i]][ingredients[j]])  
                    score += pairs[ingredients[i]][ingredients[j]]
                    num += 1
    if num > 0:
        bootstraps.append(float(score)/float(num))


import numpy

nums = sorted(bootstraps)
cdf = numpy.cumsum(nums,dtype=float)/numpy.sum(nums,dtype=float)
lastNum = nums[0]
cdfTable = {}
for i in range(len(nums)):
    if nums[i] != lastNum:
        cdfTable[int(nums[i-1])] = int(100*cdf[i-1])
cdfTable[int(nums[i])] = int(100*cdf[i])

with open('cdf.tab','w') as f:
    for key in cdfTable.keys():
        f.write("%d %d\n" % (key,cdfTable[key]))


