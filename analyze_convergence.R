library(tidyverse)

conv1 <-
  read.csv('r2d2_value_hist.csv', stringsAsFactors = FALSE, header = FALSE) %>%
  select(state = V1, time = V2, value = V3)
conv2 <-
  read.csv('termino_value_hist.csv', stringsAsFactors = FALSE, header = FALSE) %>%
  select(state = V1, time = V2, value = V3)

conv <-
  rbind(
    conv1 %>% mutate(player = "r2d2"),
    conv2 %>% mutate(player = "termino")
  )
conv$state <- as.character(conv$state)

plot <-
  ggplot(conv) +
  geom_line(aes(x = time, y = value, color = state)) +
  facet_grid(player ~ .)
print(plot)
