library(dplyr)

boardToState <- function(b, s) {
  k <- 0
  h <- 0
  for (i in 1:3) {
    for (j in 1:3) {
      b0 <- b[i, j]
      v <- case_when(b[i, j] == s ~ 0, b[i, j] == "" ~ 1, TRUE ~ 2)
      h <- h + v * 3**k
      k <- k + 1
    }
  }
  return(h)
}

stateToBoard <- function(h, s) {
  s2 <- ifelse(s == "x", "o", "x")
  b <- matrix(rep("", 9), nrow = 3, ncol = 3)
	k <- 3 * 3 - 1
	for (i in 3:1) {
		for (j in 3:1) {
			base <- 3**k
			v <- floor(h/base)
			b[i, j] <- case_when(v == 0 ~ s, v == 1 ~ "", v == 2 ~ s2)
			h <- h - v * base
			k <- k - 1
		}
	}
	return(b)
}