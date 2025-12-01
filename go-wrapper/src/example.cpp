#include "example.h"

using namespace godot;

void Example::_bind_methods() {
    ClassDB::bind_method(D_METHOD("set_message", "message"), &Example::set_message);
    ClassDB::bind_method(D_METHOD("get_message"), &Example::get_message);
    ClassDB::bind_method(D_METHOD("say_hello", "name"), &Example::say_hello);

    ADD_PROPERTY(PropertyInfo(Variant::STRING, "message"), "set_message", "get_message");
}

Example::Example() {
    // Initialize any variables here
    message = "Hello, Godot!";
}

Example::~Example() {
    // Clean up any resources here
}

void Example::set_message(const String &p_message) {
    message = p_message;
}

String Example::get_message() const {
    return message;
}

String Example::say_hello(const String &name) const {
    return "Hello, " + name + "! " + message;
}

// If you want to expose signals, you can add them like this:
// ADD_SIGNAL(MethodInfo("example_signal", PropertyInfo(Variant::STRING, "data")));
