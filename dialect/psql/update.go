package psql

import (
	"io"

	"github.com/stephenafamo/bob/builder"
	"github.com/stephenafamo/bob/expr"
	"github.com/stephenafamo/bob/mods"
	"github.com/stephenafamo/bob/query"
)

func Update(queryMods ...query.Mod[*UpdateQuery]) query.BaseQuery[*UpdateQuery] {
	q := &UpdateQuery{}
	for _, mod := range queryMods {
		mod.Apply(q)
	}

	return query.BaseQuery[*UpdateQuery]{
		Expression: q,
		Dialect:    Dialect{},
	}
}

// Trying to represent the select query structure as documented in
// https://www.postgresql.org/docs/current/sql-update.html
type UpdateQuery struct {
	expr.With
	only bool
	expr.Table
	expr.Set
	expr.FromItems
	expr.Where
	expr.Returning
}

func (u UpdateQuery) WriteSQL(w io.Writer, d query.Dialect, start int) ([]any, error) {
	var args []any

	withArgs, err := query.ExpressIf(w, d, start+len(args), u.With,
		len(u.With.CTEs) > 0, "\n", "")
	if err != nil {
		return nil, err
	}
	args = append(args, withArgs...)

	w.Write([]byte("UPDATE "))

	if u.only {
		w.Write([]byte("ONLY "))
	}

	tableArgs, err := query.ExpressIf(w, d, start+len(args), u.Table, true, "", "")
	if err != nil {
		return nil, err
	}
	args = append(args, tableArgs...)

	setArgs, err := query.ExpressIf(w, d, start+len(args), u.Set, true, " ", "")
	if err != nil {
		return nil, err
	}
	args = append(args, setArgs...)

	fromArgs, err := query.ExpressSlice(w, d, start+len(args), u.FromItems.Items, "\nFROM ", ",\n", "")
	if err != nil {
		return nil, err
	}
	args = append(args, fromArgs...)

	whereArgs, err := query.ExpressIf(w, d, start+len(args), u.Where,
		len(u.Where.Conditions) > 0, "\n", "")
	if err != nil {
		return nil, err
	}
	args = append(args, whereArgs...)

	retArgs, err := query.ExpressIf(w, d, start+len(args), u.Returning,
		len(u.Returning.Expressions) > 0, "\n", "")
	if err != nil {
		return nil, err
	}
	args = append(args, retArgs...)

	return args, nil
}

type UpdateQM struct {
	BuilderMod
	withMod[*UpdateQuery]
	mods.FromMod[*UpdateQuery]
	fromItemMod
	joinMod[*expr.FromItem]
}

func (qm UpdateQM) Only() query.Mod[*UpdateQuery] {
	return mods.QueryModFunc[*UpdateQuery](func(u *UpdateQuery) {
		u.only = true
	})
}

func (qm UpdateQM) Table(name any) query.Mod[*UpdateQuery] {
	return mods.QueryModFunc[*UpdateQuery](func(u *UpdateQuery) {
		u.Table = expr.Table{
			Expression: name,
		}
	})
}

func (qm UpdateQM) TableAs(name any, alias string) query.Mod[*UpdateQuery] {
	return mods.QueryModFunc[*UpdateQuery](func(u *UpdateQuery) {
		u.Table = expr.Table{
			Expression: name,
			Alias:      alias,
		}
	})
}

func (qm UpdateQM) Set(a, b any) query.Mod[*UpdateQuery] {
	return mods.Set[*UpdateQuery]{builder.OP("=", a, b)}
}

func (qm UpdateQM) SetArg(a, b any) query.Mod[*UpdateQuery] {
	return mods.Set[*UpdateQuery]{builder.OP("=", a, qm.Arg(b))}
}

func (qm UpdateQM) Where(e query.Expression) query.Mod[*UpdateQuery] {
	return mods.Where[*UpdateQuery]{e}
}

func (qm UpdateQM) WhereClause(clause string, args ...any) query.Mod[*UpdateQuery] {
	return mods.Where[*UpdateQuery]{qm.Statement(clause, args...)}
}

func (qm UpdateQM) Returning(expressions ...any) query.Mod[*UpdateQuery] {
	return mods.Returning[*UpdateQuery](expressions)
}
