
"""
Demonstrate a demo of caching a circuit state, acting on it, then restoring it.
"""
from circuit import *
import subprocess
import unittest

class TestGHDLAuto(unittest.TestCase):
  
  def test_compile_and_clean_testbench(self):
    before = subprocess.run(['ls'], cwd="examples/full-adder", capture_output=True, text=True)
    
    c = GHDLTestbench("examples/full-adder", ["adder"], "adder_tb", "vhd")
    c.compile()
    compiled = subprocess.run(['ls'], cwd="examples/full-adder", capture_output=True, text=True)
    
    self.assertEqual(compiled.stdout.split(), ["adder.o", "adder.vhd", "adder_tb", "adder_tb.o", "adder_tb.vhd", "e~adder_tb.o", "work-obj93.cf"])

    c.clean()
    print("done. Contents:")
    cleaned = subprocess.run(['ls'], cwd="examples/full-adder", capture_output=True, text=True)

    self.assertEqual(before.stdout, cleaned.stdout)

if __name__ == "__main__":
  unittest.main()