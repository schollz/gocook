import os

try:
    os.remove('pairing.csv')
    os.remove('pairing_foods.txt')
except:
    pass

uniqueFoods = []
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
        if "#" not in line:
            line = line.replace("_"," ")
            foods = line.split(',')
            try:
                if foods[0] not in uniqueFoods:
                    uniqueFoods.append(foods[0])
                if foods[1] not in uniqueFoods:
                    uniqueFoods.append(foods[1])
            except:
                print(foods)
            with open('pairing.csv','a') as f2:
                f2.write(line)

with open('pairing_foods.txt','w') as f:
    for food in uniqueFoods:
        f.write(food + "\n")