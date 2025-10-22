entity g_nand is
  port (A, B : in bit; Y : out bit);
end g_nand;

architecture behavior of g_nand is
begin
  Y <= not (A and B);
end behavior;