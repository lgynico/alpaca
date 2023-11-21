package template

const CSharpConfig = `
using System;

namespace ${packageName}
{
    public class ${structName}Row
    {
		${fields}
    }

    public class ${configName}Row
    {
		private
    }
}
`
