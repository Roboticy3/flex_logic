extends FlexNet

@export var indicator:PackedScene = preload("res://flex_logic/scenes/indicator.tscn")
@export var wire_renderer:PackedScene = preload("res://flex_logic/scenes/connection_renderer.tscn")
@export var overlays:Array[BaseMaterial3D] = [
	preload("res://flex_logic/materials/indicator_v0.tres"),
	preload("res://flex_logic/materials/indicator_v1.tres"),
	preload("res://flex_logic/materials/indicator_x.tres"),
	preload("res://flex_logic/materials/indicator_z.tres"),
	preload("res://flex_logic/materials/indicator_u.tres"),
]

@export var show_bits := 4

var indicator_instances:Array[GeometryInstance3D]
var wire_instances:Array[MeshInstance3D]

# Called when the node enters the scene tree for the first time.
func _ready() -> void:
	for i in min(get_size(), show_bits):
		var instance := indicator.instantiate()
		if !(instance is GeometryInstance3D):
			printerr("Low frequency visualizer cannot use ", indicator, ". Does not instantiate to GeometryInstance3D.")
			queue_free()
			return
		add_child(instance)
		create_tween().tween_property(instance, "position", Vector3.ONE * i * .5, 0.25)
		indicator_instances.push_back(instance)
	
	for c in get_connections():
		draw_line(self, c)

func _process(delta: float) -> void:
	
	#print("found logic ", get_logic(), " with nets ", get_logic().nets if get_logic() else [])
	
	for i in min(indicator_instances.size(), show_bits):
		indicator_instances[i].material_overlay = overlays[state[i]]
	
	for i in wire_instances.size():
		wire_instances[i].material_overlay = overlays[state[i]]

func draw_line(from:Node3D, to:Node3D):
	var instance := wire_renderer.instantiate()
	if !(instance is MeshInstance3D):
		printerr("Low frequency visualizer cannot use ", wire_renderer, ". Does not instanciate to MeshInstance3D.")
		queue_free()
		return
	
	wire_instances.append(instance)
	add_child(instance)
	
	var line := ImmediateMesh.new()
	instance.mesh = line
	
	line.surface_begin(Mesh.PRIMITIVE_LINES)
	line.surface_add_vertex(Vector3.ZERO)
	line.surface_add_vertex(from.to_local(to.global_position))
	line.surface_end()
	
	
