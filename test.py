import random



# ratio_fight = 0.5

# prob_fight = 0.4




a = 0.1
b = 1.3
print(a*b - a)
def run():
    e = 0.01
    av_fight = 0.0
    # av_not_fight = 0.5
    for _ in range(1000):
        # ratio_fight = random.random()
        ratio_fight = a
        prob_fight = ratio_fight * b
        fight = random.random() < prob_fight

        if fight:
            av_fight = (1-ratio_fight) * e + (1-e) * av_fight
        else:
            # print("not fight")
            av_fight = -1 * (ratio_fight) * e + (1-e) * av_fight
        # print(fight)
        # if fight:
        #     av_fight = ratio_fight * e + (1-e) * av_fight
        #     av_not_fight = (1-e) * av_not_fight
        # else:
        #     av_fight = (1-e) * av_fight
        #     av_not_fight = (1-ratio_fight) * e + (1-e) * av_not_fight
        # print(av_fight)
    return av_fight

# print(av_fight)
# print(av_not_fight)

count = 0
for i in range(1000):
    # count += 1 if run() else 0
    count += run()
print (count / 1000)
