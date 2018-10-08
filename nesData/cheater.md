
DOAE U


http://chiuinan.github.io/game/game/intro/ch/c11/sango/sango/sango.htm
6070 exp
6076 gold
61FB  Mana
008A smoke

Ji Rou , - phy damage
Who Jian ,  + fire resist
Shui Jean , + water resist
Ce Mian ,  stop enemy tac
Ji Mian , stop enemy physical att
Cheng Nei , + Castle Defenses
Yi Xin , enemy loss turn
Li Jian ,
Qi Shou , speed up
Jie Ce , remove enemy MP
An Sha,
Bei Ji , Power up ally
Fu Bing, extra attack
Tui Lu , retreat
Gui Huan , return home

-----

- 敏捷决定了攻击顺序
- General的 Region属性 决定了 这名将领 在随机遭遇中的位置。 如果你的队伍中有他，则你不会在野外随机遭遇中遇到他。
- 敌军将领的 sprite color 也会影响AP。

Locations  $0060-$0063  , Y in zone, Y zone,  X in zone, X zone
 

# Damage

## critical hits

 - the amount of damage you will do will mainly depend on your STR and A.P.
 - pick a random number between 0 and 255,
    - if the number is greater than or equal to 240 (F0), this is an excellent blow, so damage potential is 51
    - Otherwise, pick a random number between 0 and 15 , and use this number to extract the damage potential from an array
        - `20, 23, 23, 23,  25,25,25,25,25,25,25,25,   23, 23, 23,  20`
 
## Weapon Damage

```
12 = ( NONE )
25 = Dagger
38 = Flail
56 = Ax
84 = Club
128 = Spear
192 = Sabre
435 = Trident
289 = Bow
655 = Sword
1024 = Battleax
1280 = Scimitar
768 = Crossbow
1536 = Lance
1792 = Wan Sheng
1792 = Bo Ye
2048 = Qing Guang
2048 = Nu Long
2048 = Qing Long
4096 = Halberd
```

## Army size (soldiers) Strength multiplier 

```
1 to 9 –> 1
10 to 99 –> 4
100 to 999 –> 9
1000 to 9999 –> 16
10000+ –> 25
```

 - DOAE 2

```
1 to 255 –> (x / 32) + 1
256 to 2047 –> (x / 128) + 7
2048 to 8191 –> (x / 256) + 15
8192+ –> (x / 1024) + 37
```

# Tactics Damage

```
01 Lian Huo, Ye Huo, Yan Re, Da Re, Huo Shen
06 Shui Tu, Shui Xing, Shui Lei, Hong Shui, Shui Long
0b Chi Xin, Tong Xian, Yin Xian, Jin Xian, Wan Fu
10 Ji Rou, Huo Jian, Shui Jian, Ce Mian, Ji Mian
15 Cheng Nei, Yi Xin, Li Jian, Qi Shou, Jie Ce
1a An Sha, Bei Ji, Fu Bing, Tui Lu, Gui Huan
```


 - Tactic damage relies on the INT of the attacker, the INT of the opponent, and other factors.
 - each tactic has a "damage potential" associated with it

```
01 = Lian Huo = Damage = 003C = 60
02 = Ye Huo = Damage = 008C = 140
03 = Yan Re = Damage = 00AA = 170
04 = Da Re = Damage = 0514 = 1300
05 = Huo Shen = Damage = 07D0 = 2000
06 = Shui Tu = Damage = 0064 = 100
07 = Shui Xing = Damage = 00A0 = 160
08 = Shui Lei = Damage = 00C8 = 200
09 = Hong Shui = Damage = 05DC = 1500
10 = Shui Long = Damage = 09C4 = 2500
```


