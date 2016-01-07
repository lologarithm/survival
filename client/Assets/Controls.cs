using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;

namespace Assets
{
    public class ControlConst
    {
        private ControlConst(string value) { Value = value; }

        public string Value { get; set; }

        public static ControlConst Ability1 { get { return new ControlConst("Ability1"); } }
        public static ControlConst Ability2 { get { return new ControlConst("Ability2"); } }
        public static ControlConst Ability3 { get { return new ControlConst("Ability3"); } }
        public static ControlConst Ability4 { get { return new ControlConst("Ability4"); } }
        public static ControlConst Ability5 { get { return new ControlConst("Ability5"); } }
        public static ControlConst Ability6 { get { return new ControlConst("Ability5"); } }
        public static ControlConst MoveLeft { get { return new ControlConst("MoveLeft"); } }
        public static ControlConst MoveRight { get { return new ControlConst("MoveRight"); } }
        public static ControlConst MoveUp { get { return new ControlConst("MoveUp"); } }
        public static ControlConst MoveDown { get { return new ControlConst("MoveDown"); } }
    }
}
