library(tidyverse)

value_files <- Sys.glob(paste("*", "values", "csv", sep = "."))

## import state-values
memory <- c()
for (value_file in value_files) {
  name <- unlist(strsplit(value_file, split = "\\."))[1]
  memory <-
    rbind(
      memory,
      read.csv(value_file, stringsAsFactors = FALSE, header = FALSE) %>%
        select(state = V1, value = V2) %>%
        mutate(name = name)
    )
}

hist_value <-
  ggplot(memory) +
  geom_histogram(aes(x = value, fill = name),
                 alpha = 0.35, position = "identity", binwidth = 0.02) +
  labs(title = 'Histogram of state values', x = 'Value', y = 'Count')
print(hist_value)

## define each state's statistics
memory_state <-
  memory %>%
  group_by(state) %>%
  summarise(rng = max(value) - min(value),
            avg = mean(value))

hist_value_state <-
  ggplot(memory_state) +
  geom_histogram(aes(x = rng), binwidth = 0.02) +
  labs(title = 'Histogram of value range between players', x = 'Value', y = 'Count')
print(hist_value_state)
