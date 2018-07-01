library(tidyverse)

conv <-
  rbind(
    read.csv('r2d2_oldest_states_hist.csv', stringsAsFactors = FALSE, header = FALSE) %>%
      select(state = V1, time = V2, value = V3) %>%
      mutate(player = "r2d2"),
    read.csv('termino_oldest_states_hist.csv', stringsAsFactors = FALSE, header = FALSE) %>%
      select(state = V1, time = V2, value = V3) %>%
      mutate(player = "termino")
  )
conv$state <- as.character(conv$state)

plot <-
  ggplot(conv) +
  geom_line(aes(x = time, y = value, color = state), size = 1) +
  facet_grid(player ~ .)
print(plot)
