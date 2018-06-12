library(shiny)

source('board_state.R')

s <- c("x", "_", "o")

ui <- pageWithSidebar(
  headerPanel("Board encoding"),

  sidebarPanel(
    numericInput(inputId = "state",
                 label = "State code",
                 value = 0,
                 min = 0,
                 max = 19682,
                 step = 1),
    tableOutput("stateboard")
  ),

  mainPanel(
    fluidRow(
      column(3,
             h4("Column 1"),
             selectInput(inputId = "oneone",
                         label = "1-1",
                         choices = s,
                         selected = "x"),
             selectInput(inputId = "twoone",
                         label = "2-1",
                         choices = s,
                         selected = "x"),
             selectInput(inputId = "threeone",
                         label = "3-1",
                         choices = s,
                         selected = "x")
      ),
      column(3,
             h4("Column 1"),
             selectInput(inputId = "oneone",
                         label = "1-1",
                         choices = s,
                         selected = "x"),
             selectInput(inputId = "twoone",
                         label = "2-1",
                         choices = s,
                         selected = "x"),
             selectInput(inputId = "threeone",
                         label = "3-1",
                         choices = s,
                         selected = "x")
      ),
      column(3,
             h4("Column 3"),
             selectInput(inputId = "onethree",
                         label = "1-3",
                         choices = s,
                         selected = "x"),
             selectInput(inputId = "twothree",
                         label = "2-3",
                         choices = s,
                         selected = "x"),
             selectInput(inputId = "threethree",
                         label = "3-3",
                         choices = s,
                         selected = "x")
      )
    ),
    verbatimTextOutput("boardstate")
  )
)

server <- function(input, output) {
  output$stateboard <- renderTable({
    stateToBoard(as.integer(input$state), "x")
  })
  output$boardstate <- renderPrint({
    b <- matrix(data = c(input$oneone, input$onetwo, input$onethree,
                         input$twoone, input$twotwo, input$twothree,
                         input$threeone, input$threetwo, input$threethree),
                nrow = 3, ncol = 3)
    boardToState(b, "x")
  })
}

shinyApp(ui = ui, server = server)
