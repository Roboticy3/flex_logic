import subprocess

"""
Digital logic circuit represented as a GHDL testbench. 
Gates are included as separate vhdl files.
Any python GHDL solution is enough to simulate and interact with the circuit simulation.
"""
class GHDLCircuit:
  def __init__(self, include:list[str], label:str):
    self.include = include
    self.label = label
  
  def compile(self):
    subprocess.run(["ghdl", "-o", "-i", *self.include])