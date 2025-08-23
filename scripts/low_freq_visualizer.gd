extends FlexNet

@export var indicator:PackedScene = preload("res://flex_logic/scenes/indicator.tscn")
@export var overlays:Array[BaseMaterial3D] = [
	preload("res://flex_logic/materials/indicator_v0.tres"),
	preload("res://flex_logic/materials/indicator_v1.tres"),
	preload("res://flex_logic/materials/indicator_x.tres"),
	preload("res://flex_logic/materials/indicator_z.tres"),
	preload("res://flex_logic/materials/indicator_u.tres"),
]

var instances:Array[GeometryInstance3D]

# Called when the node enters the scene tree for the first time.
func _ready() -> void:
	for i in get_size():
		var instance := indicator.instantiate()
		if !(instance is GeometryInstance3D):
			printerr("Low frequency visualizer cannot use ", indicator, ". Does not instantiate to GeometryInstance3D.")
			return
		add_child(instance)
		create_tween().tween_property(instance, "position", Vector3.ONE * i, sqrt(i))
		instances.push_back(instance)

func _process(delta: float) -> void:
	var state := get_state()
	for i in instances.size():
		instances[i].material_overlay = overlays[state[i]]
