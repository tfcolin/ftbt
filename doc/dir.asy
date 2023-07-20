size (5cm, true);

settings.tex = "xelatex";
texpreamble("\usepackage{ctex}");


draw (Label ("x", EndPoint), (0, 0) -- (1, 0), E, black, Arrow);
draw (Label ("y", EndPoint), (0, 0) -- (0, -1), SE, black, Arrow);


draw (Label ("$0$", EndPoint), (2, -2) -- (3, -2), E, black, Arrow);
draw (Label ("$1$", EndPoint), (2, -2) -- (1, -2), W, black, Arrow);
draw (Label ("$2$", EndPoint), (2, -2) -- (2, -3), S, black, Arrow);
draw (Label ("$3$", EndPoint), (2, -2) -- (2, -1), N, black, Arrow);
