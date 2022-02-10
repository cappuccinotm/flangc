%{

package parser

%}

%union {}

%token QUOTE
%token SQUOTE
%token LBRACE, RBRACE
%token NUMBER
%token IDENTIFIER

%%

// expressions

Expression
  : LBRACE    Term Terms      RBRACE
  ;

Terms
  : /* empty */
  | Term Terms
  ;

Term
  : IDENTIFIER | NUMBER | Expression | List
  ;

List
  : LBRACE QUOTE Terms RBRACE
  | SQUOTE LBRACE Terms RBRACE
  ;

%%
