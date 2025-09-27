extends RichTextLabel

@export_node_path var target_path := NodePath("..") :
	set(new_target_path):
		target = get_node_or_null(new_target_path)
		target_path = new_target_path
@onready var target = get_node_or_null(target_path)

@export var prop_list:Array[StringName] = [
	&"name"
]

func _process(_delta: float) -> void:
	if !is_instance_valid(target):
		text = "No target."
		return
	text = ""
	for p in prop_list:
		if p == &"": continue
		text += "{}: {}\n\t".format([p, str(target.get(p))], "{}")
