class_name FlexNet extends Node
## A node in a digital circuit. Can "solve" by passing its value to all
## connected nodes.

## List of connections which this net should pass events to.
@export var connections:Array[FlexNet] = []

## 32-bit integer representing bus value
@export var value:int = 0
## 32-bit mask representing unset value
@export var unset:int = -1
## 32-bit mask representing conflicts
##	High value + conflict = high impedence
@export var conflict:int = 0
@export var high_impedence:int = 0

## Drive a value into each connection. 
func solve() -> Dictionary[FlexNet, float]:
	var result := {}
	for c in connections:
		pipe(c)
		result[c] = 0.0
	return result

## Drive value into a single connection. Solve for impedence rules.
## Buckwec. Table 4-2
##      0   1   x   z
##   -----------------
## 0 || 0 | x | x | 0
## 1 || x | 1 | x | 1
## x || x | x | x | x
## z || 0 | 1 | x | z
##
## All unset U are propogated.
func pipe(to:FlexNet):
	#Propogate unset
	to.unset = unset
	
	#High impedence in the bottom right corner
	to.high_impedence = high_impedence & to.high_impedence
	
	#Conflicts along the conflict column and row in the table
	#Conflicts also along conflicting values. Go figure.
	to.conflict = (conflict | to.conflict) | (value ^ to.value)
	
	#Value doesn't matter if conflict or impedence, so we can sort of 
	#"imagine" 1's in the bottom-right 3x3 of the table
	to.value = value | to.conflict | to.high_impedence
