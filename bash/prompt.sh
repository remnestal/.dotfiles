
#############################
# start of custom PS variable
Decoration1="\[\e[90m\]╔["
RegularUserPart="\[\e[1;35m\]\u"
RootUserPart="\[\e[31;5m\]\u\[\e[m\]"
Between="\[\e[90m\]@"
HostPart="\[\e[1;32m\]\h:"
PathPart="\[\e[1;34m\]\w"
Decoration2="\[\e[90m\]]\n╚>\[\e[m\]"
case `id -u` in
    0) export PS1="$Decoration1$RootUserPart$Between$HostPart$PathPart$Decoration2# ";;
    *) export PS1="$Decoration1$RegularUserPart$Between$HostPart$PathPart$Decoration2$ ";;
esac
# end of custom PS variable
#############################
