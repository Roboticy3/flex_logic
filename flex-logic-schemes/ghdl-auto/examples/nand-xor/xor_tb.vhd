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