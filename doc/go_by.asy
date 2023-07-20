size (5cm, true);

settings.tex = "xelatex";
texpreamble("\usepackage{ctex}");

pair text_box (real x0, real y0, real x1, real y1, pen p, string text) {
      draw ((x0, y0) -- (x0, y1) -- (x1, y1) -- (x1, y0) -- cycle);
      real cx = (x0 + x1) * .5;
      real cy = (y0 + y1) * .5;
      label (text, (cx, cy), p);
      return (cx, cy);
}

pair s1 = text_box (0, 0, 1, 1, red, "我方");
pair s2 = text_box (1, 1, 2, 2, blue, "敌方");
pair s3 = text_box (0, 1, 1, 2, black, "");

draw (s1 + (0, 0.3) -- s3, Arrow);
