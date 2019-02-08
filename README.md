# GoTick
Tic-tac-toe - Reinforcement learning exercise in golang

1. Introduction

The program builds a tournament of the tic-tac-toe games (https://en.wikipedia.org/wiki/Tic-tac-toe). Any number of robot and/or human players attend the tournament. In each session, two out of all players are chosen. These two players play any number of episodes. A robot has a fixed intelligence but gain experience over episodes and sessions. Each robot exports its experience to an data file which is then analyzed and visualized.

To run the program, build the executable file by `go build -o executable_file.exe` then run that executable.

2. Reinforcement learning algorithm

We use Monte-Carlo method for learning:

```
sum = 0
for t = T-1 to 0:
  sum = R[t+1] + gamma * sum
  V(x[t]) = update_func(V(x[t], sum)
end for
return V
```

where `gamma` is the discount rate of reward. The `update_func` can be chosen between the following two definitions:

(1)

```
update_func(v, s) = v + alpha * (sum - v)
```

where `alpha` is the learning rate.

(2)

```
update_func(v, s) = (n * v + sum) / (n + 1)
```

where `n` is the number of times `v` has already been updated, including the first time `v` was met.

3. Reward

Reward `R` is defined at the end of an episode, for each of the 3 outcomes: winning, losing, and draw.

4. State

At each step in an episode, the state of game for a player is defined by the game board in the player's eye; a board composed by `X`s and `O`s has to be converted to `me`s and `you`s, together with the information of who's playing the next step, to be meaningful.

reference: https://github.com/lazyprogrammer/machine_learning_examples/blob/master/rl/tic_tac_toe.py
