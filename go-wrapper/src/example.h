#pragma once

#include <godot_cpp/classes/control.hpp>
#include <godot_cpp/classes/global_constants.hpp>
#include <godot_cpp/variant/string.hpp>

namespace godot {

class Example : public Control {
    GDCLASS(Example, Control);

private:
    String message;

protected:
    static void _bind_methods();

public:
    Example();
    ~Example();

    void set_message(const String &p_message);
    String get_message() const;

    // Example method that can be called from GDScript
    String say_hello(const String &name) const;
};

}
