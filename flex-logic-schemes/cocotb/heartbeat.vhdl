library ieee;
use ieee.std_logic_1164.all;

entity heartbeat is
    port ( clk : out std_logic);
end heartbeat;

architecture behavior of heartbeat is
    constant clk_period : time := 10 ns;
begin
    clk_process: process
    begin --processes loop, I didn't notice before.
        clk <= '0'; -- send 0 to the clk port
        wait for clk_period/2;
        clk <= '1';
        wait for clk_period/2;
    end process;
end behavior;               
