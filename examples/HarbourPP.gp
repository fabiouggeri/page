grammar HarbourPP;

/*********************************************************************************
                              RULES
**********************************************************************************/

PreProcessor : Statements EOI;

/*********************************************************************************
                              LEXER
**********************************************************************************/

@Fragment
MultilineDoubleQuoteStringSegment : ('\n' | '\"' | BreakString)!+;

@Fragment
MultilineSingleQuoteStringSegment : ('\n' | '\'' | BreakString)!+;

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

@Fragment
OneSpace : (' ' | '\t' | '\f');

NewLine : '\n' | '\r\n';

AloneLineComment : '*' ('\n' | EOI)!*;

Whitespace : OneSpace+;

BlockComment : '/*' ('*'! | ('*' '/'!))* '*/';

LineComment : ('//' | '&&') ('\n' | EOI)!*;

BracketString : '[' ('\n' | ']')!* ']';

LogicalLiteral : ".T." | ".F." | ".Y." | ".N.";

DoubleQuoteString : '"' ('\n' | '\"')!* '"';

SingleQuoteString : '\'' ('\n' | '\'')!* '\'';

// @Atomic
DateTimeLiteral : ('0d' Digit Digit Digit Digit Digit Digit Digit Digit) 
                | ("d" ( ('"' DatePattern ("T"? TimePattern)? '"') 
					        | ('\'' DatePattern ("T"? TimePattern)? '\'') 
							  | ('[' DatePattern ("T"? TimePattern)? ']'))) 
							  | ('{' '^' (IntegerNumber ('-' | '/') IntegerNumber ('-' | '/') IntegerNumber ','?)? TimePattern? Ignore '}');

NumberLiteral : (IntegerNumber '.' IntegerNumber) | (IntegerNumber '.' Letter!) | IntegerNumber | ('.' IntegerNumber);

MultiLineDoubleQuoteString : '"' (MultilineDoubleQuoteStringSegment BreakString)+ ('\n' | '\"')!* '"';

MultiLineSingleQuoteString : '\'' (MultilineSingleQuoteStringSegment BreakString)+ ('\n' | '\'')!* '\'';

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
Statements : Statement*;

@SkipNode
Ignore : Whitespace | BlockComment | LineComment | ContinueNL;

@Memoize
ContinueNL : OneSpace* ';' (BlockComment | LineComment | Whitespace)* NewLine;

@Name(Spacing)
@SkipNode
OptionalSpacing : Ignore*;

@SkipNode
Spacing : Ignore+;

@SkipNode
Statement : DirectiveStatement | EmptyStatement | AnyStatement;

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

EmptyStatement : Ignore* AloneLineComment?;

AnyStatement : AnyRules;

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
OptionalMatchMarkerChunk : ResultSep! (OptionalMatchMarker | MatchMarker | Identifier | Literal | Ignore | EscapedChar | Separator);

@SkipNode
MatchMarker : IdMarker | ListMarker | RestrictMarker | WildMarker | ExtendedMarker | IdentifierMarker;

@SkipNode
Literal : LogicalLiteral | DoubleQuoteString | SingleQuoteString | DateTimeLiteral | NumberLiteral | MultiLineDoubleQuoteString | MultiLineSingleQuoteString;

@SkipNode
MatchChunk : ResultSep! (OptionalMatchMarker | MatchMarker | Identifier | Literal | Ignore | EscapedChar | ']' | Separator);

IdMarker : '<' Identifier '>';

ListMarker : '<' Identifier ',' '...' '>';

RestrictMarker : '<' Identifier ':' RestrictValues '>';

WildMarker : '<' '*' Identifier '*' '>';

ExtendedMarker : '<' '(' Identifier ')' '>';

IdentifierMarker : '<' '!' Identifier '!' '>';

RestrictValues : RestrictValue (',' RestrictValue)*;

@SkipNode
RestrictValue : (('>' | ',')! (Identifier | Literal | Ignore | EscapedChar | Separator))+;

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
OptionalResultMarkerChunk : OptionalResultMarker | ResultMarker | Identifier | Literal | Ignore | EscapedChar | Separator;

@SkipNode
ResultChunk : OptionalResultMarker | ResultMarker | Identifier | Literal | Ignore | EscapedChar | ']' | Separator;

BeginDumpBlock : "pragma" "begindump";

EndDumpBlock : "pragma" "enddump";

@SkipNode
BracketSequence : ('[' (']'! AnyRule)+ ']') | BracketString;

@SkipNode
AnyRule : BracketSequence | Literal | Ignore | Identifier | Separator;

BreakString : ';' OneSpace* NewLine;
