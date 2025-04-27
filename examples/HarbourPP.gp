grammar HarbourPP;

/*********************************************************************************
                              MAIN RULE
**********************************************************************************/

PreProcessor : Statements EOI;

/*********************************************************************************
                              LEXER
**********************************************************************************/
@Fragment
Letter : [a-z] | [A-Z];

@Fragment
Digit : [0-9];

@Fragment
IntegerNumber : Digit+;

@Fragment
TimePattern : IntegerNumber (':' IntegerNumber (':' IntegerNumber)? ('.' IntegerNumber)?)? ("AM" | "PM")?;

@Fragment
DatePattern : IntegerNumber ('-' | '/' | '.') IntegerNumber ('-' | '/' | '.') IntegerNumber;

NewLine : '\n' | '\r\n';

@Ignore
Whitespace : (' ' | '\t' | '\f')+;

@Ignore
BlockComment : '/*' ('*'! | ('*' '/'!))* '*/';

@Ignore
LineComment : ('//' | '&&') ('\n' | EOI)!*;

BracketString : '[' ('\n' | ']')!* ']';

LogicalLiteral : ".T." | ".F." | ".Y." | ".N.";

DoubleQuoteString : '"' ('\n' | '"')!* '"';

SingleQuoteString : '\'' ('\n' | '\'')!* '\'';

DateTime : ('0d' Digit Digit Digit Digit Digit Digit Digit Digit) 
         | ('d"' ('\n' | '"')!* '"') 
			| ("d'" ('\n' | '\'')!* "'") 
			| ('d[' ('\n' | ']')!* ']');

// @Atomic
DateTimeLiteral : DateTime 
                | ('{' '^' (IntegerNumber ('-' | '/') IntegerNumber ('-' | '/') IntegerNumber ','?)? TimePattern? '}');

NumberLiteral : (IntegerNumber '.' IntegerNumber) | (IntegerNumber '.' Letter!) | IntegerNumber | ('.' IntegerNumber);

Identifier : ([A-Z] | [a-z] | '_') ([A-Z] | [a-z] | [0-9] | '_')*;

////////////////////  KEYWORDS //////////////////////////
Separator : ".or." | ".and." | ".not." | ':=' | '==' | '!=' | '>=' | '<=' | '->' | '++' | '--' | '+=' | '-=' 
          | '*=' | '/=' | '%=' | '^=' | '**' | '^^' | '<<' | '>>' | '::' | '<>' | '...' | '&&' | '||' | '^^' 
          | '**=' | '$' | ',' | '>' | '+' | '*' | '!' | '-' | '/' | '(' | ')' | ':' | '{' | '}' | '%' | '\\' 
          | '?' | '~' | '.' | '@' | '|' | '&' | '=' | '^' | '#' | ';' | '<' | '[' | '"' | '\'' | '_';

IfDef : "ifdef":4;

IfNDef : "ifndef":4;

ElseIf : "elseif":5;

EndIf : "endif":4;

Undef : "undef":4;

Include : "include":4;

Define : "define":4;

StdOut : "stdout":4;

Command : "command":4;

YCommand : "ycommand":4;

Uncommand : "uncommand":4;

Xuncommand : "xuncommand":4;

Yuncommand : "yuncommand":4;

XCommand : "xcommand":4;

Untranslate : "untranslate":4;

Xuntranslate : "xuntranslate":4;

Yuntranslate : "yuntranslate":4;

Translate : "translate":4;

YTranslate : "ytranslate":4;

XTranslate : "xtranslate":4;

/*********************************************************************************
                              PARSER
**********************************************************************************/
@SkipNode
Statements : (Statement) (NewLine Statement?)*;

@SkipNode
Statement : DirectiveStatement | AloneLineComment | AnyStatement;

EndStmt : NewLine | EOI;

@SkipNode
DirectiveStatement : '#' 
                     ( DefineDirective 
                     | StdOutDirective 
                     | CommandDirective 
                     | XCommandDirective 
                     | YCommandDirective 
                     | TranslateDirective 
                     | XTranslateDirective 
                     | YTranslateDirective 
                     | IfDefDirective 
                     | IfNDefDirective 
                     | ElseDirective 
                     | ElseIfDirective 
                     | EndIfDirective 
                     | UndefDirective 
                     | ErrorDirective 
                     | IncludeDirective 
                     | LineDirective 
                     | UncommandDirective 
                     | XUncommandDirective 
                     | YUncommandDirective 
                     | UntranslateDirective 
                     | XUntranslateDirective 
                     | YUntranslateDirective 
                     | DumpBlock);

// EmptyStatement : NewLine;

AnyStatement : AnyRules;

@Ignore
AloneLineComment : '*' AnyRule*;

@SkipNode
AnyRules : AnyRule+;

DefineDirective : Define Identifier DefineParameters? ResultRules?;

StdOutDirective : StdOut ResultRules;

CommandDirective : Command MatchPattern ResultSep ResultPattern?;

XCommandDirective : XCommand directivePattern;

YCommandDirective : YCommand MatchPattern ResultSep ResultPattern?;

TranslateDirective : Translate directivePattern;

XTranslateDirective : XTranslate directivePattern;

YTranslateDirective : YTranslate directivePattern;

IfDefDirective : IfDef Identifier;

IfNDefDirective : IfNDef Identifier;

ElseDirective : "else";

ElseIfDirective : ElseIf Identifier;

EndIfDirective : EndIf;

UndefDirective : Undef Identifier;

ErrorDirective : "error" ResultRules;

IncludeName : (DoubleQuoteString | SingleQuoteString);

IncludeDirective : Include IncludeName;

LineDirective : "line";

UncommandDirective : Uncommand undefDirectivePattern;

XUncommandDirective : Xuncommand undefDirectivePattern;

YUncommandDirective : Yuncommand undefDirectivePattern;

UntranslateDirective : Untranslate undefDirectivePattern;

XUntranslateDirective : Xuntranslate undefDirectivePattern;

YUntranslateDirective : Yuntranslate undefDirectivePattern;

DumpBlock : BeginDumpBlock (EndDumpBlock | EOI)!* EndDumpBlock;

DefineParameters : '(' ParametersList? ')';

ResultRules : AnyRules;

@SkipNode
ParametersList : Identifier (',' Identifier)*;

MatchPattern : MatchChunk+;

// @Atomic
ResultSep : '=' '>';

ResultPattern : ResultChunk+;

@SkipNode
directivePattern : MatchPattern ResultSep ResultPattern?;

@SkipNode
undefDirectivePattern : MatchPattern;

@Atomic
EscapedChar : '\\' '\n'!;

OptionalMatchMarker : '[' OptionalMatchMarkerPattern ']';

@Name(MatchPattern)
OptionalMatchMarkerPattern : OptionalMatchMarkerChunk+;

@SkipNode
OptionalMatchMarkerChunk : ResultSep! (OptionalMatchMarker | MatchMarker | Identifier | Literal | EscapedChar | Separator);

@SkipNode
MatchMarker : IdMarker | ListMarker | RestrictMarker | WildMarker | ExtendedMarker | IdentifierMarker;

@SkipNode
Literal : LogicalLiteral | DoubleQuoteString | SingleQuoteString | DateTimeLiteral | NumberLiteral;

@SkipNode
MatchChunk : ResultSep! (OptionalMatchMarker | MatchMarker | Identifier | Literal | EscapedChar | ']' | Separator);

IdMarker : '<' Identifier '>';

ListMarker : '<' Identifier ',' '...' '>';

RestrictMarker : '<' Identifier ':' RestrictValues '>';

WildMarker : '<' '*' Identifier '*' '>';

ExtendedMarker : '<' '(' Identifier ')' '>';

IdentifierMarker : '<' '!' Identifier '!' '>';

RestrictValues : RestrictValue (',' RestrictValue)*;

@SkipNode
RestrictValue : (('>' | ',')! (Identifier | Literal | EscapedChar | Separator))+;

NullMarker : '<' '-' Identifier '-' '>';

@SkipNode
ResultMarker : IdMarker | DumbStringifyMarker | NormalStringifyMarker | SmartStringifyMarker | BlockifyMarker | LogifyMarker | NullMarker;

DumbStringifyMarker : '#' '<' Identifier '>';

NormalStringifyMarker : '<' (('"' Identifier '"') | ('\'' Identifier '\'')) '>';

SmartStringifyMarker : '<' '(' Identifier ')' '>';

BlockifyMarker : '<' '{' Identifier '}' '>';

LogifyMarker : '<' '.' Identifier '.' '>';

OptionalResultMarker : '[' OptionalResultMarkerPattern ']';

@Name(ResultPattern)
OptionalResultMarkerPattern : OptionalResultMarkerChunk+;

@SkipNode
OptionalResultMarkerChunk : OptionalResultMarker | ResultMarker | Identifier | Literal | EscapedChar | Separator;

@SkipNode
ResultChunk : OptionalResultMarker | ResultMarker | Identifier | Literal | EscapedChar | ']' | Separator;

BeginDumpBlock : "pragma" "begindump";

EndDumpBlock : "pragma" "enddump";

@SkipNode
BracketSequence : ('[' (']'! AnyRule)+ ']') | BracketString;

@SkipNode
AnyRule : BracketSequence | Literal | Identifier | Separator;

