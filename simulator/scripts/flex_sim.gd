extends Node
## An event stack for a full circuit. Can be bootstrapped by passing events
## manually. Continues until either no events are left or events have too much
## delay to process on the current step.

## Event """"queue""""
var events:Dictionary[FlexNet, float]

## Solve the event queue over a certain amount of time.
func solve(delta:float):
	var working := true
	while working:
		var events_keys = events.keys()
		working = false
		var new_events := {}
		
		# In theory, this inner loop should be "easy" to parallelize.
		# This would require placing a lock on values changed by each `solve`
		#	call. 
		# A proper priority queue system would also make this better.
		for e in events_keys:
			if events[e] <= delta:
				working = true
				new_events[e] = e.solve()
		
		for e in events_keys:
			if e in new_events:
				events.erase(e)
				events.merge(new_events[e])
			else:
				events[e] -= delta
