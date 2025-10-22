import subprocess, os, glob
from typing import Iterable

"""
Digital logic circuit represented as a GHDL testbench. 
Gates are included as separate vhdl files.
Any python GHDL solution is enough to simulate and interact with the circuit simulation.
"""
class GHDLCircuit:
  pass

"""
Simplified version of GHDLCircuit renders a static testbench as a circuit.
"""
class GHDLTestbench:
  def __init__(self, target_dir:str, components:list[str], testbench:str, ext:str) -> None:
    self.target_dir = target_dir
    self.components = components
    self.testbench = testbench
    self.ext = ext
  
  def get_component_paths(self) -> Iterable[str]:
    return map(lambda c: f"{c}.{self.ext}", self.components)
  
  def get_testbench_path(self) -> str:
    return f"{self.testbench}.{self.ext}"
  
  def compile(self) -> None:
    analysis_args = ["ghdl", "-a", 
      *self.get_component_paths(), self.get_testbench_path()
    ]
    subprocess.run(analysis_args, cwd=self.target_dir)
    
    elaboration_args = ["ghdl", "-e", self.testbench]
    subprocess.run(elaboration_args, cwd=self.target_dir)
  
  def clean(self):
    subprocess.run(["rm", 
      *glob.glob("*.o", root_dir=self.target_dir), 
      *glob.glob("*.cf", root_dir=self.target_dir), 
      self.testbench
    ], cwd=self.target_dir)

"""
Next step up from GHDLTestbench, allowing for construction of the testbench
from text.
"""
class GHDLStringTestbench(GHDLTestbench):
  def __init__(self, target_dir, components, testbench, ext, testbench_text):
    super().__init__(target_dir, components, testbench, ext)
    self.testbench_text = testbench_text
  
  def compile(self):
    with open(os.path.join(self.target_dir, self.get_testbench_path()), 'w') as f:
      f.write(self.testbench_text)
    
    return super().compile()

  def clean(self):
    os.remove(os.path.join(self.target_dir, self.get_testbench_path()))

    return super().clean()