library(tidyverse)

memory1 <-
  read.csv('r2d2_values.csv', stringsAsFactors = FALSE, header = FALSE) %>%
  select(state = V1, value1 = V2)
memory2 <-
  read.csv('termino_values.csv', stringsAsFactors = FALSE, header = FALSE) %>%
  select(state = V1, value2 = V2)

memory <-
  inner_join(memory1,
             memory2,
             by = "state") %>%
  mutate(sum = value1 + value2)

hist_value <-
  ggplot(memory) +
  geom_histogram(aes(x = value1), binwidth = 0.05, fill = 'darkred', alpha = 0.3) +
  geom_histogram(aes(x = value2), binwidth = 0.05, fill = 'darkgreen', alpha = 0.3) +
  labs(title = 'Histogram of state values (red = p1, green = p2)', x = 'Value', y = 'Count')
print(hist_value)

hist_value_sum <-
  ggplot(memory) +
  geom_histogram(aes(x = sum), binwidth = 0.001) +
  labs(title = 'Histogram of (p1 value + p2 value)', x = 'Value', y = 'Count')
print(hist_value_sum)

c <- cor(x = memory$value1, y = memory$value2)
print(c)