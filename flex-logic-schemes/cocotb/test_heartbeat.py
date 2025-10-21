import cocotb
from cocotb.triggers import Timer
from cocotb.simtime import get_sim_time

@cocotb.test()
async def heartbeat_monitor(dut):
    while True:
        print("%d: dut clk %s \t\t\t (click to continue)" % (get_sim_time('ns'), dut.clk.value))
        input()
        await Timer(5, unit="ns")
