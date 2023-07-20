size (5cm, true);

settings.tex = "xelatex";
texpreamble("\usepackage{ctex}");

int i, j, k;

string [] text = {"D", "F", "S", "DX", "P", "A"};
pen [] tcolor = {red, blue, black};

for (i = 0; i < 3; ++ i) {
      tcolor[i] += linewidth(1);
}

for (i = 0; i < 2; ++ i) {
      for (j = 0; j < 3; ++ j) {
            draw (box ((j * 2, i * 2), ((j + 1) * 2, (i + 1) * 2)));
            if (i != 1 || j != 0) 
                  draw ((j * 2 + 1, i * 2) -- (j * 2 + 1, (i + 1) * 2), dashed);
            for (k = 0; k < 2; ++ k) {
                  if (i != 1 || j != 0) 
                        label (text[i * 3 + j], (j * 2 + k + 0.5, i * 2 + 1), tcolor[k]);
            }
            label (text[3], (1, 3), tcolor[2]);
      }
}
