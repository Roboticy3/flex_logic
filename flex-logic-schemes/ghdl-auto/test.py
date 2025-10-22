
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
    
    self.assertTrue(set(compiled.stdout.split()).issuperset(["adder.o", "adder.vhd", "adder_tb", "adder_tb.o", "adder_tb.vhd", "e~adder_tb.o"]))

    c.clean()
    print("done. Contents:")
    cleaned = subprocess.run(['ls'], cwd="examples/full-adder", capture_output=True, text=True)

    self.assertEqual(before.stdout, cleaned.stdout)
  
  def test_compile_and_clean_custom(self):
    before = subprocess.run(['ls'], cwd="examples/nand-xor", capture_output=True, text=True)
    
    c = GHDLTestbench("examples/nand-xor", ["g_nand"], "xor_tb", "vhd")
    c.compile()
    compiled = subprocess.run(['ls'], cwd="examples/nand-xor", capture_output=True, text=True)
    
    self.assertTrue(set(compiled.stdout.split()).issuperset(["g_nand.o", "g_nand.vhd", "xor_tb", "xor_tb.o", "xor_tb.vhd", "e~xor_tb.o"]))
    print("compiled contents:", compiled.stdout.split())

    c.clean()
    print("done. Contents:")
    cleaned = subprocess.run(['ls'], cwd="examples/nand-xor", capture_output=True, text=True)

    self.assertEqual(before.stdout, cleaned.stdout)
  
  def test_compile_string_testbench(self):
    design = """
-- Declaration
entity xor_tb is
end xor_tb;

architecture behavior of xor_tb is
  -- Declarations of gates that will be used.
  component g_nand
    port (A, B : in bit; Y : out bit);
  end component;

  -- Declarations of gate instances that will be connected.
  for nand_0: g_nand use entity work.g_nand;
  for nand_1: g_nand use entity work.g_nand;
  for nand_2: g_nand use entity work.g_nand;
  for nand_3: g_nand use entity work.g_nand;

  -- Declarations of inputs and outputs.
  signal A, B, Y : bit;

  -- Declarations of internal nets
  signal n1, n2, n3 : bit;
begin
  -- Declaration of gate connections.. Definitely an adjustment
  nand_0: g_nand port map (A => A, B => B, Y => n1);
  nand_1: g_nand port map (A => A, B => n1, Y => n2);
  nand_2: g_nand port map (A => n1, B => B, Y => n3);
  nand_3: g_nand port map (A => n2, B => n3, Y => Y);

  -- Test process
  process
    type pattern_type is record
      -- Inputs
      A, B : bit;
      -- Expected outputs
      Y : bit;
    end record;
    type pattern_array is array (natural range <>) of pattern_type;
    constant patterns : pattern_array :=
      (('0', '0', '1'),
       ('0', '1', '1'),
       ('1', '0', '1'),
       ('1', '1', '0'));
  begin
    -- Check each pattern.
    for i in patterns'range loop
      A <= patterns(i).A;
      B <= patterns(i).B;
      wait for 10 ns;
      assert Y = patterns(i).Y
        report "bad output value" severity error;
    end loop;
    assert false report "end of test" severity note;
    wait;
  end process;
end behavior;
"""

    cwd = "examples/nand-xor-generated"

    before = subprocess.run(['ls'], cwd=cwd, capture_output=True, text=True)
    
    c = GHDLStringTestbench(cwd, ["g_nand"], "xor_tb", "vhd", design)
    c.compile()
    compiled = subprocess.run(['ls'], cwd=cwd, capture_output=True, text=True)
    
    self.assertTrue(set(compiled.stdout.split()).issuperset(["g_nand.o", "g_nand.vhd", "xor_tb", "xor_tb.o", "xor_tb.vhd", "e~xor_tb.o"]))
    print("compiled contents:", compiled.stdout.split())

    c.clean()
    print("done. Contents:")
    cleaned = subprocess.run(['ls'], cwd=cwd, capture_output=True, text=True)

    self.assertEqual(before.stdout, cleaned.stdout)

if __name__ == "__main__":
  unittest.main()