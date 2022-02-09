%{

package parser

%}

%union {}

%token LESS GREATER EQUAL NOT_EQUAL
%token PLUS MINUS MULTIPLY DIVIDE
%token NUMBER
%token IDENTIFIER
%token NULL
%token LBRACE, RBRACE

%%

CompoundName
  : IDENTIFIER
  ;

LeftPart
  : CompoundName
  ;

// expressions

Relation
  : Expression
  | Expression RelationalOperator Expression
  ;

RelationalOperator
  : LESS | GREATER | EQUAL | NOT_EQUAL
  ;

Expression
  :           Term Terms
  | LBRACE    Term Terms      RBRACE
  | LBRACE AddSign Term Terms RBRACE
  ;

AddSign
  : PLUS | MINUS
  ;

Terms
  : /* empty */
  | AddSign Term Terms
  ;

Term
  : Factor Factors
  ;

Factors
  : /* empty */
  | MultSign Factor Factors
  ;

MultSign
  : MULTIPLY | DIVIDE
  ;

Factor
  : NUMBER
  | LeftPart
  | NULL
  ;

%%
